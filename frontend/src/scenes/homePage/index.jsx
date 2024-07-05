import Author from "components/Author";
import Books from "components/Books";
import Search from "components/Search";
import User from "components/User";
import { useSelector } from "react-redux";

const HomePage = () => {
  const user = useSelector((state) => state.user);
  const token = useSelector((state) => state.token);

  return (
    <div className="grid grid-cols-3 w-full">
      <div className="flex m-4 w-full justify-center">
        {user === null ? <User /> : <Author />}
      </div>
      <div className="flex m-4 justify-center p-4 max-h-screen overflow-y-auto">
        <Books />
      </div>
      <div className="flex m-4 justify-center p-4">
        <Search />
      </div>
    </div>
  );
};

export default HomePage;
