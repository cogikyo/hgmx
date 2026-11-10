import { animate, τmax } from "../motion.js";

export default function (el) {
  const face = el.querySelector("g");
  animate(face, { rotate: [0, -45, 15, 0] }, { duration: τmax });
}
