import { animate, τmax } from "../motion.js";

export default function (el) {
  const chilli = el.querySelector("g");
  animate(chilli, { rotate: [0, -10, 6, -3, 0] }, { duration: τmax });
}
