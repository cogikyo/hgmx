import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const pointer = el.querySelector('path[d^="M9.037"]');
  animate(pointer, { translateX: [0, -2, 1, 0], translateY: [0, -2, 1, 0] }, { duration: τmid });

  const clicks = ["m6 12", "m5.1 8", "M7.2 2.2", "M14 4.1"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(clicks, { scale: [0, 1.75, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.03, { startDelay: τmin }) });
}
