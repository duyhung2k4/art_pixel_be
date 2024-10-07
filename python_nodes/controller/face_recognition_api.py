from flask import Blueprint, request, jsonify
import json
import face_recognition
import numpy as np
import os

face_recognition_bp = Blueprint('face_recognition', __name__)

@face_recognition_bp.route('/recognize_faces', methods=['POST'])
def recognize_faces_from_db():
    data = request.json

    try:
        faces = data["faces"]
        known_face_encodings = [np.array(face["faceEncoding"]) for face in faces]
        known_profile_ids = [face["profileId"] for face in faces]

        input_image_path = data["input_image_path"]
        if not os.path.exists(input_image_path):
            return jsonify({"result": "-2", "message": "Input image path does not exist."}), 400

    except KeyError as e:
        return jsonify({"result": "-2", "error": f"Missing key: {str(e)}"}), 400
    except Exception as e:
        return jsonify({"result": "-2", "error": str(e)}), 400

    try:
        def recognize_face_in_image(input_image_path):
            image_to_check = face_recognition.load_image_file(input_image_path)

            # Lấy vị trí khuôn mặt
            face_locations = face_recognition.face_locations(image_to_check)

            # Lấy vị trí các đặc trưng khuôn mặt (lông mày, mắt, mũi, miệng, ...)
            face_landmarks_list = face_recognition.face_landmarks(image_to_check)

            new_face_locations = []

            # Lặp qua từng khuôn mặt
            for i, (top, right, bottom, left) in enumerate(face_locations):
                # Lấy vị trí lông mày
                if len(face_landmarks_list) > i:
                    landmarks = face_landmarks_list[i]

                    # Lấy vị trí lông mày
                    left_eyebrow = landmarks.get("left_eyebrow", [])
                    right_eyebrow = landmarks.get("right_eyebrow", [])

                    if left_eyebrow and right_eyebrow:
                        # Tính tọa độ trung bình của lông mày
                        eyebrow_top_y = min([point[1] for point in left_eyebrow + right_eyebrow])
                        
                        # Giảm vùng khuôn mặt để chỉ lấy từ trên lông mày đổ xuống
                        new_top = eyebrow_top_y  # Vị trí từ trên lông mày
                        new_face_locations.append((new_top, right, bottom, left))

            # Lấy mã hóa khuôn mặt trong vùng đã điều chỉnh
            face_encodings = face_recognition.face_encodings(image_to_check, new_face_locations)

            if len(face_encodings) == 0:
                return -3  # Không tìm thấy khuôn mặt

            for face_encoding in face_encodings:
                matches = face_recognition.compare_faces(known_face_encodings, face_encoding)

                if not any(matches):
                    continue  # Không tìm thấy khớp nào, tiếp tục với mã hóa tiếp theo

                face_distances = face_recognition.face_distance(known_face_encodings, face_encoding)
                best_match_index = np.argmin(face_distances)

                if matches[best_match_index]:
                    profile_id = known_profile_ids[best_match_index]
                    accuracy = 1 - face_distances[best_match_index]  # Tính toán độ chính xác

                    print(f"ProfileID:{profile_id} / result: {round(accuracy * 100, 2)}")

                    if round(accuracy * 100, 2) >= 70.00:
                        return profile_id
                    return -3  # Trả về ID của profile

            return -3  # Không tìm thấy khớp nào

        profile_id = recognize_face_in_image(input_image_path)
        if profile_id == -3:
            return jsonify({"result": "-3", "message": "No matching faces found."})
        
        return jsonify({"result": str(profile_id)})

    except Exception as e:
        return jsonify({"result": "-4", "error": str(e)}), 500
