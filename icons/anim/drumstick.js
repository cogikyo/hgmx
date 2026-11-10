import { animate, τmax } from "../motion.js";

export default function (el) {
  const drumstick = el.querySelector("g");
  animate(drumstick, { rotate: [0, -12, 8, 0] }, { duration: τmax });
}
