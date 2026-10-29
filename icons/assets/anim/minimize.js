import { animate, τmax } from "../motion.js";

export default function (el) {
  const corners = el.querySelector("g");
  animate(corners, { scale: [1.2, 0.9, 1] }, { duration: τmax });
}
