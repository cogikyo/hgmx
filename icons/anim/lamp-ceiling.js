import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const lamp = el.querySelector("g");
  animate(lamp, { translateY: [0, 1.5, 0] }, { duration: τmin });

  const shade = el.querySelector("g g");
  animate(shade, { scale: [1, 0.85, 1.15, 1] }, { duration: τmid, delay: τmin });
}
