import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const lines = ["M3 5h18", "M3 12h3.458", "M3 19h3.832"].map((d) => el.querySelector(`path[d="${d}"]`));
  animate(lines, { translateX: [-3, 0], opacity: [0, 1] }, { duration: τmin, delay: stagger(0.06) });

  const hands = el.querySelector('path[d^="M16 13"]');
  animate(hands, { rotate: [0, 360] }, { duration: τmid, delay: τmin });
}
