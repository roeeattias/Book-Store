import Author from "components/Author";
import Books from "components/Books";
import User from "components/User";
import { useSelector } from "react-redux";

const HomePage = () => {
  const user = useSelector((state) => state.user);
  const token = useSelector((state) => state.token);

  return (
    <div className="grid grid-cols-3 w-full">
      <div className="flex m-4 w-full justify-center">
        {user === null ? <User></User> : <Author></Author>}
      </div>
      <div className="flex m-4 justify-center p-4">
        <Books />
      </div>
    </div>
  );
};

export default HomePage;
