import { animate, τmax } from "../motion.js";

export default function (el) {
  const fan = el.querySelector("g");
  animate(fan, { translateY: [0, -2, 0], rotate: [0, 8, -3, 0] }, { duration: τmax });
}
