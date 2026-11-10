import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const pie = el.querySelectorAll('path[d^="M21 16"], path[d^="M2 16"], path[d^="m8.5"], path[d^="M15.5"]');
  animate(pie, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const steam = ["M7 2", "M12 2", "M17 2"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(steam, { translateY: [3, -1, 0], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06, { startDelay: τmin }) });
}
