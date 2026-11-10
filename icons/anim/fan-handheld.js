import { animate, τmax } from "../motion.js";

export default function (el) {
  const fan = el.querySelector("g");
  animate(fan, { rotate: [0, -12, 8, -4, 0] }, { duration: τmax });
}
