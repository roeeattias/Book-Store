import Books from "components/Books";

const HomePage = () => {
  return (
    <div className="grid grid-cols-3 w-full">
      <div className="bg-slate-600 m-4"></div>
      <div className=" bg-background flex m-4 justify-center p-4">
        <Books />
      </div>
    </div>
  );
};

export default HomePage;
