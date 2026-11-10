import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const axes = el.querySelector('path[d^="M3 3"]');
  animate(axes, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const [falling, rising] = el.querySelectorAll("g");
  animate(falling, { translateY: [0, 1.5, -0.5, 0] }, { duration: τmid, delay: τmin });
  animate(rising, { translateY: [0, -1.5, 0.5, 0] }, { duration: τmid, delay: τmid });
}
