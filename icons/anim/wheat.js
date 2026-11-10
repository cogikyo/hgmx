import { animate, τmax } from "../motion.js";

export default function (el) {
  const wheat = el.querySelector("g");
  animate(wheat, { rotate: [0, -6, 4, -2, 0] }, { duration: τmax });
}
