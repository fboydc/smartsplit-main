import { useContext } from "react";
import Context from "../../Context";

const AuthError = () => {

   const { authError } = useContext(Context);
   let message = "";
   if (authError.error_code === "401") {
      message = "Invalid username or password combination.";
   }

   if (authError.error_code === "500") {
      //return <div>System has encountered an error. Please try again later.</div>
      message = "System has encountered an error. Please try again later.";
   }    

   if (authError.error_code === "504")
      message = "System appears to be down. Please try again later."
   return (
        <div>{message}</div>
   )
}

export default AuthError;