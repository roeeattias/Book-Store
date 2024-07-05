import { useSelector } from "react-redux";
import { redirect, useNavigate } from "react-router-dom";

const User = () => {
  const booksBought = useSelector((state) => state.bought);
  const navigate = useNavigate();

  return (
    <div className="m-5 w-full font-inika flex flex-col gap-4">
      <div className="border-black rounded-lg border-2 flex flex-row items-center pt-3 pb-3 pl-4 gap-4">
        <img
          src={`${process.env.PUBLIC_URL}/images/profile_picture.jpg`}
          alt="Book"
          className="w-16 h-16 object-cover rounded-full"
        />
        <div className="flex flex-col gap-1">
          <div className="text-xl font-semibold">Anonymous</div>
          <div className="border border-gray-950" />
          <div className="text-sm text-opacity-40 text-black font-semibold">
            Books bought: {booksBought}
          </div>
        </div>
      </div>
      <button
        className="bg-BecomeAnAuthorButton rounded-lg p-3 text-white font-semibold w-full text-xl"
        onClick={() => {
          navigate("/login");
        }}
      >
        Become An Author
      </button>
    </div>
  );
};

export default User;
