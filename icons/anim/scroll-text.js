import { animate, stagger, τmax } from "../motion.js";

export default function (el) {
  const paths = el.querySelectorAll("path");

  const scrollPaths = [paths[2], paths[3]];
  animate(scrollPaths, { y: [-10, 2, 0] }, { duration: τmax });
  animate(scrollPaths, { pathLength: [0, 1] }, { duration: τmax });

  const textPaths = [paths[1], paths[0]];
  animate(
    textPaths,
    { opacity: [0, 1], x: [-3, 0], y: [-2, 0] },
    { duration: 0.3, delay: stagger(0.12) },
  );
}
