from flask import Blueprint, request, jsonify
import cv2
import face_recognition
import os

face_detection_bp = Blueprint('face_detection', __name__)

@face_detection_bp.route('/detect_single_face', methods=['POST'])
def detect_single_face():
    image_path = request.json.get("input_image_path")
    
    if not image_path or not os.path.exists(image_path):
        return jsonify({"result": False, "error": "Image path is invalid."}), 400

    try:
        image = cv2.imread(image_path)
        if image is None:
            return jsonify({"result": False, "error": "Unable to load image."}), 400
        
        # Resize image to improve speed (if image is large)
        if image.shape[0] > 800 or image.shape[1] > 800:
            scale_factor = 800 / max(image.shape[0], image.shape[1])
            image = cv2.resize(image, (0, 0), fx=scale_factor, fy=scale_factor)

        face_locations = face_recognition.face_locations(image)

        if len(face_locations) == 1:
            return jsonify({"result": True})
        elif len(face_locations) == 0:
            return jsonify({"result": False, "message": "No face detected."}), 200
        else:
            return jsonify({"result": False, "message": "Multiple faces detected."}), 200

    except Exception as e:
        return jsonify({"result": False, "error": str(e)}), 500
