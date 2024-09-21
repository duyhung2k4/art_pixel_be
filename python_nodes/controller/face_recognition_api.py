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
            return jsonify({"result": "-2"}), 400
    except Exception as e:
        return jsonify({"result": "-2", "error": str(e)}), 400

    try:
        def recognize_face_in_image(input_image_path):
            image_to_check = face_recognition.load_image_file(input_image_path)
            face_locations = face_recognition.face_locations(image_to_check)
            face_encodings = face_recognition.face_encodings(image_to_check, face_locations)

            # Nếu không phát hiện khuôn mặt nào
            if len(face_encodings) == 0:
                return "-3"  # Trả về -3 nếu không có khuôn mặt nào

            for face_encoding in face_encodings:
                matches = face_recognition.compare_faces(known_face_encodings, face_encoding)
                
                # Nếu không có mặt nào khớp
                if not any(matches):
                    return "-3"  # Trả về -3 nếu không có mặt nào khớp

                face_distances = face_recognition.face_distance(known_face_encodings, face_encoding)
                best_match_index = np.argmin(face_distances)

                if matches[best_match_index]:
                    profile_id = known_profile_ids[best_match_index]

                    if len(face_locations) == 1:
                        return f"{profile_id}"

            return "-3"  # Trả về -3 nếu không tìm thấy khớp

        message = recognize_face_in_image(input_image_path)
        return jsonify({"result": message})
    except Exception as e:
        return jsonify({"result": "-4", "error": str(e)}), 500
