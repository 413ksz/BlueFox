
const Technologies = () => {
    const technologies = [
        { name: "SolidJS", url: "https://www.solidjs.com/" },
        { name: "Tailwind CSS", url: "https://tailwindcss.com/" },
    ];
  return (
    <section id="technologies" className="bg-gray-950/50 py-16 md:py-24">
    <div className="container mx-auto px-4">
        <h2 className="text-3xl sm:text-4xl font-semibold text-center mb-12 text-white">
            Technologies
        </h2>
        <div className="flex flex-wrap justify-center gap-6">
            {technologies.map((tech, index) => (
                <a
                    key={index}
                    href={tech.url}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="bg-gray-800/80 backdrop-blur-md rounded-xl px-4 py-2 shadow-md border border-gray-700 text-gray-300 hover:text-white hover:bg-blue-500/20 transition-colors"
                    aria-label={`Learn more about ${tech.name}`}
                >
                    {tech.name}
                </a>
            ))}
        </div>
    </div>
</section>
  )
}

export default Technologies
