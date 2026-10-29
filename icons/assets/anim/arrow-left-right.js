import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const left = el.querySelector("g.left");
  animate(left, { translateX: [0, -4, 1, 0] }, { duration: τmid });

  const right = el.querySelector("g.right");
  animate(right, { translateX: [0, 4, -1, 0] }, { duration: τmid, delay: τmin });
}
