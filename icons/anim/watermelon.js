import { animate, τmid } from "../motion.js";

export default function (el) {
  const slice = el.querySelector("g");
  animate(slice, { rotate: [0, 8, -3, 0] }, { duration: τmid });
}
