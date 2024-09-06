import ImageUpload from "./ImageUpload";
import { GiCancel } from "react-icons/gi";

const PublishBookModal = ({
  setTitle,
  setAuthor,
  setQuantity,
  setPrice,
  setProfilePicture,
  profilePicture,
  closeModal,
}) => {
  return (
    <div className="flex flex-col border border-black rounded-md p-5 justify-around gap-4">
      <div className="flex flex-row gap-10 font-bold">
        <div className="flex flex-col w-full gap-4">
          <div className="titleOfBook">
            <div>Title Of Book</div>
            <form className="w-full">
              <input
                type="text"
                placeholder="title"
                onChange={(e) => {
                  setTitle(e.target.value);
                }}
                className="w-full border border-black rounded-sm border-opacity-20 bg-searchBoxFill bg-opacity-20 text-xs font-bold pl-2 pt-2 pb-2"
              />
            </form>
          </div>
          <div className="author">
            <div>Author</div>
            <form className="w-full">
              <input
                type="text"
                placeholder="author"
                onChange={(e) => {
                  setAuthor(e.target.value);
                }}
                className="w-full border border-black rounded-sm border-opacity-20 bg-searchBoxFill bg-opacity-20 text-xs font-bold pl-2 pt-2 pb-2"
              />
            </form>
          </div>
        </div>
        <div className="flex flex-col w-full gap-4">
          <div className="quantity">
            <div>Quantity</div>
            <form className="w-full">
              <input
                type="text"
                placeholder="quantity"
                onChange={(e) => {
                  setQuantity(e.target.value);
                }}
                className="w-full border border-black rounded-sm border-opacity-20 bg-searchBoxFill bg-opacity-20 text-xs font-bold pl-2 pt-2 pb-2"
              />
            </form>
          </div>
          <div className="price">
            <div>Price</div>
            <form className="w-full">
              <input
                type="text"
                placeholder="price"
                onChange={(e) => {
                  setPrice(e.target.value);
                }}
                className="w-full border border-black rounded-sm border-opacity-20 bg-searchBoxFill bg-opacity-20 text-xs font-bold pl-2 pt-2 pb-2"
              />
            </form>
          </div>
        </div>
      </div>
      <div className="w-full flex flex-row justify-between items-center">
        <ImageUpload
          setProfilePicture={setProfilePicture}
          profilePicture={profilePicture}
        />
        <GiCancel
          className="h-5 w-5 hover:text-red-600 hover:scale-105 transition ease-in-out duration-300"
          onClick={closeModal}
        />
      </div>
    </div>
  );
};

export default PublishBookModal;
