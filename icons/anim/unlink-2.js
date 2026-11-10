import { animate, τmax } from "../motion.js";

export default function (el) {
  const link = el.querySelector("path");
  animate(link, { scaleX: [1, 1.2, 0.95, 1] }, { duration: τmax });
}
