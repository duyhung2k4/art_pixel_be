from flask import Blueprint, request, jsonify
import cv2
import os
import dlib
import numpy as np

calculate_head_pose_bp = Blueprint('calculate_head_pose', __name__)

# Hàm tính toán góc xoay khuôn mặt
def calculate_head_pose(landmarks):
    image_points = np.array([
        (landmarks.part(30).x, landmarks.part(30).y),  # Mũi
        (landmarks.part(8).x, landmarks.part(8).y),    # Cằm
        (landmarks.part(36).x, landmarks.part(36).y),  # Điểm mắt trái
        (landmarks.part(45).x, landmarks.part(45).y),  # Điểm mắt phải
        (landmarks.part(48).x, landmarks.part(48).y),  # Điểm miệng trái
        (landmarks.part(54).x, landmarks.part(54).y)   # Điểm miệng phải
    ], dtype="double")

    size = (640, 480)
    focal_length = size[1]
    center = (size[1] // 2, size[0] // 2)

    model_points = np.array([
        (0.0, 0.0, 0.0),            # Mũi
        (0.0, -330.0, -65.0),       # Cằm
        (-225.0, 170.0, -135.0),    # Điểm mắt trái
        (225.0, 170.0, -135.0),     # Điểm mắt phải
        (-150.0, -150.0, -125.0),   # Điểm miệng trái
        (150.0, -150.0, -125.0)     # Điểm miệng phải
    ])

    camera_matrix = np.array(
        [[focal_length, 0, center[0]],
         [0, focal_length, center[1]],
         [0, 0, 1]], dtype="double"
    )

    dist_coeffs = np.zeros((4, 1))

    success, rotation_vector, translation_vector = cv2.solvePnP(
        model_points, image_points, camera_matrix, dist_coeffs)

    rotation_matrix, _ = cv2.Rodrigues(rotation_vector)

    roll = np.arctan2(rotation_matrix[2, 1], rotation_matrix[2, 2])
    pitch = np.arctan2(-rotation_matrix[2, 0], np.sqrt(rotation_matrix[2, 1]**2 + rotation_matrix[2, 2]**2))
    yaw = np.arctan2(rotation_matrix[1, 0], rotation_matrix[0, 0])

    return np.degrees(np.array([roll, pitch, yaw]))

@calculate_head_pose_bp.route('/calculate_head_pose', methods=['POST'])
def calculate_pose():
    try:
        # Lấy đường dẫn ảnh từ yêu cầu POST
        data = request.json
        image_path = data['input_image_path']

        # Đọc và xử lý ảnh
        image = cv2.imread(image_path)
        if image is None:
            return jsonify({"error": "Image not found"}), 404

        gray = cv2.cvtColor(image, cv2.COLOR_BGR2GRAY)

        # Phát hiện khuôn mặt
        detector = dlib.get_frontal_face_detector()
        current_dir = os.path.dirname(os.path.abspath(__file__))
        predictor_path = os.path.join(current_dir, "shape_predictor_68_face_landmarks.dat")
        predictor = dlib.shape_predictor(predictor_path)
        
        faces = detector(gray)

        if len(faces) == 0:
            return jsonify({"error": "No face detected"}), 400

        for face in faces:
            landmarks = predictor(gray, face)
            angles = calculate_head_pose(landmarks)

            yaw = angles[2]
            
            if yaw > 7.5:
                return jsonify({ "result": False })
            if yaw < -7.5:
                return jsonify({ "result": False })

            return jsonify({ "result": True })

    except Exception as e:
        print(e)
        return jsonify({"result": False, "error": str(e)}), 500
