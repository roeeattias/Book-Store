import { useState } from "react";
import { useDispatch } from "react-redux";
import { useNavigate } from "react-router-dom";
import { setLogin } from "state";

const LoginPage = () => {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const navigate = useNavigate();
  const dispatch = useDispatch();

  const onSubmit = async () => {
    if (username && password) {
      const response = await fetch("http://localhost:8080/login", {
        method: "POST",
        credentials: "include",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ username: username, password: password }),
      });

      if (response.status === 200) {
        const user = await response.json();
        dispatch(setLogin({ user: user }));
        navigate("/");
      } else {
        alert("Username or password are incorrect");
      }
    }
    if (username === "") {
      alert("Please enter your username");
    }
    if (password === "") {
      alert("Please enter your password");
    }
  };

  return (
    <div className="flex flex-col items-center justify-between w-full pt-48 font-inika font-bold">
      <div className="flex flex-col items-center">
        <div className="text-5xl mb-4">Login To Your Account</div>
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
        </div>
        <button
          onClick={onSubmit}
          className="bg-becomeAnAuthorButton bg-opacity-90 w-3/4 rounded-lg p-4 text-white font-semibold text-base hover:bg-opacity-100"
        >
          Login To My Account
        </button>
        <p className="text-opacity-40 text-black mt-2 text-sm">
          Don't Have An Account Already ?{" "}
          <a
            href="/signup"
            className="text-becomeAnAuthorButton text-opacity-90 hover:text-opacity-100"
          >
            Signup
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

export default LoginPage;
