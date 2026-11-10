import { animate, τmax } from "../motion.js";

export default function (el) {
  const ear = el.querySelector("g");
  animate(ear, { rotate: [0, -8, 5, 0] }, { duration: τmax });
}
