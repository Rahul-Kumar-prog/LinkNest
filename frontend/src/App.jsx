import { BrowserRouter, Route, Routes } from "react-router-dom";
import Loginpage from "./Components/Login-Signup/Loginpage";
import Signuppage from "./Components/Login-Signup/Signuppage";
import Home from "./pages/Home";
import ChannelWindow from "./pages/ChannelWindow";

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/channel/:channelId" element={<ChannelWindow />} />
        <Route path="/Loginpage" element={<Loginpage />} />
        <Route path="/login" element={<Loginpage />} />
        <Route path="/Signuppage" element={<Signuppage />} />
        <Route path="/signup" element={<Signuppage />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
