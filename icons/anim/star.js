import { animate, τmax, τmid } from "../motion.js";

export default function (el) {
  const star = el.querySelector("path");
  animate(star, { scale: [0, 1.25, 1], rotate: [0, 360] }, { duration: τmax });
  animate(star, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });
}
