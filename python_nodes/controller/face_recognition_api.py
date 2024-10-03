from flask import Blueprint, request, jsonify
from deepface import DeepFace
import os
import numpy as np

face_recognition_bp = Blueprint('face_recognition', __name__)

@face_recognition_bp.route('/recognize_faces', methods=['POST'])
def recognize_faces_from_db():
    data = request.json

    try:
        # Nhận face encodings từ JSON input
        faces = data["faces"]
        known_face_encodings = [np.array(face["faceEncoding"]) for face in faces]
        known_profile_ids = [face["profileId"] for face in faces]

        # Kiểm tra đường dẫn ảnh đầu vào
        input_image_path = data["input_image_path"]
        if not os.path.exists(input_image_path):
            return jsonify({"result": "-2"}), 400
    except Exception as e:
        return jsonify({"result": "-2", "error": str(e)}), 400

    try:
        def recognize_face_in_image(input_image_path):
            # Sử dụng DeepFace để trích xuất các đặc trưng khuôn mặt (embeddings)
            face_encodings = DeepFace.represent(img_path=input_image_path, model_name="VGG-Face", enforce_detection=False)
            
            if len(face_encodings) == 0:
                return "-3", 0.0  # Không có khuôn mặt nào được phát hiện

            best_match_profile_id = None
            best_match_distance = float('inf')  # Khởi tạo với giá trị lớn nhất

            # So sánh với face_encodings đã lưu trong database
            for face_encoding_dict in face_encodings:
                face_encoding = np.array(face_encoding_dict['embedding'])  # Trích xuất giá trị embedding
                
                # Tính khoảng cách giữa các vector face_encodings
                distances = [np.linalg.norm(face_encoding - known_encoding) for known_encoding in known_face_encodings]
                min_distance = min(distances)
                best_match_index = distances.index(min_distance)

                if min_distance < 2:  # Ngưỡng khoảng cách
                    profile_id = known_profile_ids[best_match_index]
                    confidence = 1 - (min_distance / 2)  # Độ chính xác, có thể thay đổi công thức tùy thuộc vào yêu cầu
                    
                    # In độ chính xác ra console
                    print(f"Profile ID: {profile_id}, Confidence: {confidence:.2f}")
                    
                    if confidence > 0.7:
                        return f"{profile_id}"
                    return f"-3"

            return "-3"  # Không tìm thấy khuôn mặt nào khớp

        message = recognize_face_in_image(input_image_path)
        return jsonify({"result": message})
    
    except Exception as e:
        print(e)
        return jsonify({"result": "-4", "error": str(e)}), 500
