import ImageUpload from "components/ImageUpload";
import { useState } from "react";
import { useDispatch } from "react-redux";
import { useNavigate } from "react-router-dom";
import { setLogin } from "state";

const SignUpPage = () => {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [profilePicture, setProfilePicture] = useState(null);
  const navigate = useNavigate();
  const dispatch = useDispatch();

  const onSubmit = async () => {
    if (username && password && profilePicture) {
      const response = await fetch("http://localhost:8080/signup", {
        method: "POST",
        credentials: "include",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          username: username,
          password: password,
          image_url: profilePicture,
        }),
      });

      if (response.status === 201) {
        const user = await response.json();
        dispatch(setLogin({ user: user }));
        setProfilePicture(null);
        navigate("/");
      } else {
        alert("Username already exists");
      }
    }
    if (username === "") {
      alert("Please enter your username");
    } else if (password === "") {
      alert("Please enter your password");
    } else if (profilePicture === null) {
      alert("Please upload profile picture");
    }
  };

  return (
    <div className="flex flex-col items-center justify-between w-full pt-48 font-inika font-bold">
      <div className="flex flex-col items-center">
        <div className="text-5xl mb-4">Become An Author</div>
        <p className="text-opacity-60 text-black mb-6">
          “Unleash your creativity, publish your story. Explore, buy, and
          connect with your favorite authors.”
        </p>
        <div className="w-2/3 flex flex-col gap-3 mb-5">
          <form className="w-full flex justify-center items-center">
            <input
              type="text"
              placeholder="Username"
              onChange={(e) => {
                setUsername(e.target.value);
              }}
              className="w-full border border-black rounded-md p-3 border-opacity-20 bg-searchBoxFill bg-opacity-20 text-sm font-bold pl-5"
            />
          </form>
          <form className="w-full flex justify-center items-center">
            <input
              type="text"
              placeholder="Password"
              onChange={(e) => {
                setPassword(e.target.value);
              }}
              className="w-full border border-black rounded-md p-3 border-opacity-20 bg-searchBoxFill bg-opacity-20 text-sm font-bold pl-5"
            />
          </form>
          <ImageUpload
            setProfilePicture={setProfilePicture}
            profilePicture={profilePicture}
          />
        </div>
        <button
          onClick={onSubmit}
          className="bg-becomeAnAuthorButton bg-opacity-90 w-3/4 rounded-lg p-4 text-white font-semibold text-base hover:bg-opacity-100"
        >
          Create My Author Account
        </button>
        <p className="text-opacity-40 text-black mt-2 text-sm">
          Have An Account Already ?{" "}
          <a
            href="/login"
            className="text-becomeAnAuthorButton text-opacity-90 hover:text-opacity-100"
          >
            Login
          </a>
        </p>
      </div>
      <p className="flex items-center justify-center w-full self-end mb-4 gap-1 font-bold text-opacity-40 text-black text-sm">
        By continuing, you agree to our{" "}
        <a className="text-becomeAnAuthorButton text-opacity-90">
          Terms of Service
        </a>{" "}
        and confirm that you have read our{" "}
        <a className="text-becomeAnAuthorButton text-opacity-90">
          Privacy Policy.
        </a>
      </p>
    </div>
  );
};

export default SignUpPage;
