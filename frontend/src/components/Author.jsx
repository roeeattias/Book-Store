import { useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import PublishBookModal from "./PublishBookModal";
import { useNavigate } from "react-router-dom";
import { addBook, setLogout } from "state";
import { TbLogout2 } from "react-icons/tb";

const Author = ({ setVisitedAuthor }) => {
  const user = useSelector((state) => state.user);
  const navigate = useNavigate();
  const dispatch = useDispatch();
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [title, setTitle] = useState("");
  const [author, setAuthor] = useState("");
  const [quantity, setQuantity] = useState("");
  const [price, setPrice] = useState("");
  const [profilePicture, setProfilePicture] = useState(null);

  const openModal = () => setIsModalOpen(true);
  const closeModal = () => setIsModalOpen(false);

  const publishBook = async () => {
    if (!isModalOpen) {
      openModal();
    } else {
      if (!(title && author && quantity && price && profilePicture)) {
        alert("Please enter all the required fields");
      } else {
        try {
          const response = await fetch("http://localhost:8080/publishBook", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            credentials: "include",
            body: JSON.stringify({
              title: title,
              author: author,
              quantity: parseInt(quantity),
              price: parseInt(price),
              image_url: profilePicture,
            }),
          });
          if (response.status === 201) {
            const book = await response.json();
            dispatch(addBook({ book: book }));
            closeModal();
          } else if (response.status === 511) {
            navigate("/login");
          } else {
            alert("Error publishing book");
          }
        } catch {
          alert("Caanot publish book without internet connection");
        }
      }
    }
  };

  return (
    <div className="m-5 w-full font-inika flex flex-col gap-3">
      <div className="border-black rounded-lg border-2 flex flex-row items-center pt-3 pb-3 pl-4 gap-4 justify-between">
        <div
          className="flex flex-row w-full items-center gap-4 hover:scale-105 hover:translate-x-1 transition ease-in-out delay-50 duration-300"
          onClick={() => {
            setVisitedAuthor(user);
          }}
        >
          <img
            src={user.image_url}
            alt="Book"
            className="w-16 h-16 object-cover rounded-full"
          />
          <div className="flex flex-col gap-1">
            <div className="text-xl font-semibold">{user.username}</div>
            <div className="border border-gray-950" />
            <div className="text-sm text-opacity-40 text-black font-semibold">
              Published: {user.publishedBooks.length} books
            </div>
            <div className="text-sm text-opacity-40 text-black font-semibold">
              Sold: 600 books
            </div>
          </div>
        </div>
        <TbLogout2
          className="h-8 w-8 mr-5 hover:text-becomeAnAuthorButton"
          onClick={() => {
            dispatch(setLogout());
          }}
        />
      </div>
      {isModalOpen && (
        <PublishBookModal
          setTitle={setTitle}
          setAuthor={setAuthor}
          setPrice={setPrice}
          setQuantity={setQuantity}
          setProfilePicture={setProfilePicture}
          profilePicture={profilePicture}
          closeModal={closeModal}
        />
      )}
      <button
        className="bg-becomeAnAuthorButton rounded-lg p-3 text-white font-semibold w-full text-xl"
        onClick={publishBook}
      >
        Publish New Book
      </button>
    </div>
  );
};

export default Author;
