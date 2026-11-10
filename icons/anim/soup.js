import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const bowl = el.querySelectorAll('path[d^="M4 12"], path[d^="M7 21"]');
  animate(bowl, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const steam = ["M6.252", "M11.248", "M16.252"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(steam, { translateY: [3, -1, 0], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06, { startDelay: τmin }) });
}
