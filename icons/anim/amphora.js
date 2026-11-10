import { animate, τmax } from "../motion.js";

export default function (el) {
  const amphora = el.querySelector("g");
  animate(amphora, { rotate: [0, -8, 6, -3, 0] }, { duration: τmax });
}
