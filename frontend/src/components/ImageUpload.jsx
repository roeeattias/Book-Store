import React from "react";

function ImageUpload({ setProfilePicture, profilePicture }) {
  const handleFileChange = (event) => {
    const file = event.target.files[0];
    if (file) {
      const reader = new FileReader();
      reader.onloadend = () => {
        setProfilePicture({
          file: file,
          dataUrl: reader.result,
        });
      };
      reader.readAsDataURL(file);
    }
  };

  return (
    <div className="flex justify-center items-center text-becomeAnAuthorButton text-opacity-90 text-sm">
      <input
        type="file"
        accept="image/*"
        onChange={handleFileChange}
        style={{ display: "none" }}
        id="uploadInput"
      />
      <label htmlFor="uploadInput">
        <span>
          Upload Image {profilePicture ? profilePicture.file.name : ""}
        </span>{" "}
      </label>
    </div>
  );
}

export default ImageUpload;
