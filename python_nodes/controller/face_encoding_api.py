from flask import Blueprint, request, jsonify
from deepface import DeepFace
import os

face_encoding_bp = Blueprint('face_encoding', __name__)

@face_encoding_bp.route('/face_encoding', methods=['POST'])
def face_encoding():
    # Nhận đường dẫn thư mục từ JSON yêu cầu POST
    data = request.get_json()
    directory_path = data.get("directory_path")

    if not directory_path or not os.path.exists(directory_path):
        return jsonify({"result": "error", "message": "Directory path is invalid."}), 400

    list_face_encoding = []

    try:
        # Duyệt qua từng file ảnh trong thư mục
        for image_file in os.listdir(directory_path):
            image_path = os.path.join(directory_path, image_file)

            # Mã hóa khuôn mặt sử dụng DeepFace
            embedding = DeepFace.represent(img_path=image_path, model_name="VGG-Face", enforce_detection=False)

            if embedding:
                # Nếu có nhiều embedding, đưa tất cả vào list_face_encoding
                for emb in embedding:
                    list_face_encoding.append(emb['embedding'])

        # Trả về danh sách mã hóa khuôn mặt dưới dạng JSON
        return jsonify({"result": "success", "face_encodings": list_face_encoding})
    
    except Exception as e:
        # Bắt lỗi và trả về thông báo lỗi
        return jsonify({"result": "error", "message": str(e)}), 500
