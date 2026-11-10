import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const lamp = el.querySelector("g");
  animate(lamp, { scaleY: [1, 0.85, 1] }, { duration: τmin });

  const shade = el.querySelector("path[transform-origin]");
  animate(shade, { scale: [1, 0.85, 1.15, 1] }, { duration: τmid, delay: τmin });
}
