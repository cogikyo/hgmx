import { animate, τmax } from "../motion.js";

export default function (el) {
  const dress = el.querySelector("g");
  animate(dress, { rotate: [0, -10, 6, -3, 0] }, { duration: τmax });
}
