import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const [left, right, top] = el.querySelectorAll("g");
  animate(left, { translateX: [-4, 1, 0], opacity: [0, 1] }, { duration: τmid });
  animate(right, { translateX: [4, -1, 0], opacity: [0, 1] }, { duration: τmid });
  animate(top, { translateY: [-6, 1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
