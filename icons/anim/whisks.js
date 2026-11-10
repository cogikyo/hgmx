import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const [left, right] = el.querySelectorAll("g");
  animate(left, { rotate: [0, 10, -4, 0] }, { duration: τmax });
  animate(right, { rotate: [0, -10, 4, 0] }, { duration: τmid, delay: τmin });
}
