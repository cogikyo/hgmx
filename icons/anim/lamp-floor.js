import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const lamp = el.querySelector("g");
  animate(lamp, { scaleY: [1, 0.9, 1] }, { duration: τmin });

  const shade = el.querySelector('path[d^="M17.929"]');
  animate(shade, { scale: [1, 0.85, 1.15, 1] }, { duration: τmid, delay: τmin });
}
