import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const stroller = el.querySelectorAll("path, circle");
  animate(stroller, { translateX: [0, 3, -1, 0] }, { duration: τmax });

  const head = el.querySelector('circle[cx="11"]');
  animate(head, { translateY: [1, -2, 0] }, { duration: τmin, delay: τmid });
}
