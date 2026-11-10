import { animate, τmax } from "../motion.js";

export default function (el) {
  const mic = el.querySelector("g");
  animate(mic, { rotate: [0, -12, 8, 0] }, { duration: τmax });
}
