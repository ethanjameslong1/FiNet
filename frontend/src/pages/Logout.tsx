import { Button } from "../components/Button";

const Logout = () => {
  return (
    <div>
      <Button to="/logout" type="submit">
        Logout
      </Button>
    </div>
  );
};

export default Logout;
