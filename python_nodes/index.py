import sys
from flask import Flask, jsonify
from controller.face_recognition_api import face_recognition_bp
from controller.face_detection_api import face_detection_bp
from controller.face_encoding_api import face_encoding_bp
from controller.calculate_head_pose_api import calculate_head_pose_bp

app = Flask(__name__)

# Đăng ký các blueprint
app.register_blueprint(face_recognition_bp)
app.register_blueprint(face_detection_bp)
app.register_blueprint(face_encoding_bp)
app.register_blueprint(calculate_head_pose_bp)

@app.route('/ping', methods=['GET'])
def ping():
    return jsonify({"message": "pong"}), 200

if __name__ == "__main__":
    # Lấy port từ argument dòng lệnh (nếu có)
    if len(sys.argv) > 1:
        port = int(sys.argv[1])
    else:
        port = 5000  # Giá trị mặc định nếu không truyền port

    print(f"Flask server is running on port {port}")
    app.run(host='0.0.0.0', port=port)
