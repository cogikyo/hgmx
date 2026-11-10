import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const rabbit = el.querySelector("g");
  animate(rabbit, { translateY: [0, -3, 0.5, 0] }, { duration: τmid });

  const ear = el.querySelector('path[d^="M20 8.54"]');
  animate(ear, { rotate: [0, 15, -5, 0] }, { duration: τmid, delay: τmin });
}
