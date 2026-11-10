import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const cord = el.querySelector('path[d="M4 3v14"]');
  animate(cord, { scaleY: [1, 1.1, 1] }, { duration: τmid });

  const pull = el.querySelector("circle");
  animate(pull, { translateY: [0, 1.4, 0] }, { duration: τmid });

  const slats = ["M20 7H8", "M20 11H8", "M8 15h12", "M10 19h10"].map((d) => el.querySelector(`path[d="${d}"]`));
  animate(slats, { translateY: [0, 1.5, 0] }, { duration: 0.3, delay: stagger(0.03, { startDelay: τmin }) });
}
