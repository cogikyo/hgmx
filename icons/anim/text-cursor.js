import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const caret = el.querySelector("g");
  animate(caret, { scaleY: [0, 1.15, 1] }, { duration: τmid });
  animate(caret, { opacity: [1, 0, 1] }, { duration: τmin, delay: τmid });
}
