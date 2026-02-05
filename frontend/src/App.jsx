import {BrowserRouter, Routes, Route} from "react-router-dom";
import Home from "./pages/Home";
import Loginpage from "./Components/Login-Signup/Loginpage";
import Signuppage from "./Components/Login-Signup/Signuppage";

function App(){
  return (
    <BrowserRouter>
    <Routes>
      <Route path="/" element={<Home/>} />
      <Route path="/Loginpage" element={<Loginpage/>} />
      <Route path="/Signuppage" element={<Signuppage/>} />
    </Routes>
    </BrowserRouter>
  );
}

export default App;