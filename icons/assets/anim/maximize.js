import { animate, τmax } from "../motion.js";

export default function (el) {
  const corners = el.querySelector("g");
  animate(corners, { scale: [0.8, 1.1, 1] }, { duration: τmax });
}
