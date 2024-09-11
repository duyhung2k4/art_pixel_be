import face_recognition
import sys
import os

def face_encoding(directory_path):
    list_face_encoding = []
    for image_file in os.listdir(directory_path):
        image_path = os.path.join(directory_path, image_file)
        
        new_image = face_recognition.load_image_file(image_path)
        new_face_encodings = face_recognition.face_encodings(new_image)
        
        if len(new_face_encodings) > 0:
            new_face_encoding = new_face_encodings[0]
            list_face_encoding.append(new_face_encoding)
        else:
            print(f"No faces found in {image_file}")
    
    return list_face_encoding

if __name__ == "__main__":
    directory_path = sys.argv[1]
    result = face_encoding(directory_path)
    
    if result:
        encoding_str = f"{[encoding.tolist() for encoding in result]}"
        print(encoding_str)
    else:
        print("No face encodings were found.")
