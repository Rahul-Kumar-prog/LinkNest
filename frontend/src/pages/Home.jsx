import Navbar from "../Components/Landing/Navbar";
import Hero from "../Components/Landing/Hero";
import Sidebar from "../Components/Sidebar/Sidebar";

export default function Home(){
    return (
        <div className="min-h-screen bg-gray-900">
            <Sidebar />
            <div className="ml-64">
                <Navbar />
                <Hero />
            </div>
        </div>
    );
}

