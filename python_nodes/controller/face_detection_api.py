from flask import Blueprint, request, jsonify
from deepface import DeepFace
import os
import cv2

face_detection_bp = Blueprint('face_detection', __name__)

@face_detection_bp.route('/detect_single_face', methods=['POST'])
def detect_single_face():
    image_path = request.json.get("input_image_path")

    if not image_path or not os.path.exists(image_path):
        return jsonify({"result": False, "error": "Image path is invalid."}), 400

    try:
        # Dùng OpenCV để load hình ảnh
        image = cv2.imread(image_path)

        # Phát hiện khuôn mặt bằng DeepFace
        face_objs = DeepFace.detectFace(img_path=image_path, detector_backend='opencv', enforce_detection=False)
        
        is_single_face = len(face_objs) == 1
        
        if is_single_face:
            return jsonify({"result": True})
        return jsonify({"result": False})
    
    except Exception as e:
        return jsonify({"result": False, "error": str(e)}), 500
