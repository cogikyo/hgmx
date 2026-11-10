import { animate, τmax } from "../motion.js";

export default function (el) {
  const band = el.querySelector("path");
  animate(band, { scaleX: [1, 0.8, 1.1, 1] }, { duration: τmax });
}
