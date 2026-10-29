import { animate, τmax, τmid } from "../motion.js";

export default function (el) {
  const band = el.querySelector("path");
  animate(band, { scaleX: [1, 0.8, 1.1, 1] }, { duration: τmax });
  animate(band, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
}
